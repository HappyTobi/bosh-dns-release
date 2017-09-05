package records

import (
	"encoding/json"
	"net"
	"strings"

	"errors"

	"strconv"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/miekg/dns"
)

type RecordSet struct {
	Domains []string
	Records []Record
}

func (r RecordSet) Resolve(fqdn string) ([]string, error) {
	if net.ParseIP(fqdn) != nil {
		return []string{fqdn}, nil
	}

	var ips []string

	if strings.HasPrefix(fqdn, "q-") {
		matcher := strings.SplitN(fqdn, ".", 2)
		if len(matcher) < 2 {
			return ips, errors.New("domain is malformed")
		}
		encodedQuery := strings.TrimPrefix(matcher[0], "q-")
		filter, err := parseCriteria(encodedQuery)
		if err != nil {
			return ips, err
		}

		for _, record := range r.Records {
			recordName := record.Fqdn(false)
			if recordName == matcher[1] && filter.isAllowed(record) {
				ips = append(ips, record.IP)
			}
		}
	} else {
		for _, record := range r.Records {
			compare := record.Fqdn(true)
			if compare == fqdn {
				ips = append(ips, record.IP)
			}
		}
	}

	return ips, nil
}

func CreateFromJSON(j []byte, logger boshlog.Logger) (RecordSet, error) {
	s := RecordSet{}
	swap := struct {
		Keys  []string        `json:"record_keys"`
		Infos [][]interface{} `json:"record_infos"`
	}{}

	err := json.Unmarshal(j, &swap)
	if err != nil {
		return RecordSet{}, err
	}

	s.Records = make([]Record, 0, len(swap.Infos))
	s.Domains = []string{}

	idIndex := -1
	groupIndex := -1
	networkIndex := -1
	deploymentIndex := -1
	ipIndex := -1
	domainIndex := -1
	azIDIndex := -1
	instanceIndexIndex := -1

	for i, k := range swap.Keys {
		switch k {
		case "id":
			idIndex = i
		case "instance_group":
			groupIndex = i
		case "network":
			networkIndex = i
		case "deployment":
			deploymentIndex = i
		case "ip":
			ipIndex = i
		case "domain":
			domainIndex = i
		case "az_id":
			azIDIndex = i
		case "instance_index":
			instanceIndexIndex = i
		default:
			continue
		}
	}

	domains := map[string]struct{}{}
	countKeys := len(swap.Keys)

	for index, info := range swap.Infos {
		countInfo := len(info)
		if countInfo != countKeys {
			logger.Warn("RecordSet", "Unbalanced records structure. Found %d fields of an expected %d at record #%d", countInfo, countKeys, index)
			continue
		}

		var domainIndexStr string
		if !assertStringValue(&domainIndexStr, info, domainIndex, "domain", index, logger) {
			continue
		}

		domain := dns.Fqdn(domainIndexStr)
		domains[domain] = struct{}{}

		record := Record{Domain: domain}

		if !assertStringValue(&record.ID, info, idIndex, "id", index, logger) {
			continue
		} else if !assertStringValue(&record.Group, info, groupIndex, "group", index, logger) {
			continue
		} else if !assertStringValue(&record.Network, info, networkIndex, "network", index, logger) {
			continue
		} else if !assertStringValue(&record.Deployment, info, deploymentIndex, "deployment", index, logger) {
			continue
		} else if !assertStringValue(&record.IP, info, ipIndex, "ip", index, logger) {
			continue
		}

		assertStringValue(&record.AZID, info, azIDIndex, "az_id", index, logger)
		assertStringIntegerValue(&record.InstanceIndex, info, instanceIndexIndex, "instance_index", index, logger)

		s.Records = append(s.Records, record)
	}

	for domain := range domains {
		s.Domains = append(s.Domains, domain)
	}

	return s, nil
}

func assertStringIntegerValue(field *string, info []interface{}, fieldIdx int, fieldName string, infoIdx int, logger boshlog.Logger) bool {
	if fieldIdx < 0 {
		return false
	}

	float64Value, ok := info[fieldIdx].(float64) // golang default type for numeric fields
	if !ok {
		logger.Warn("RecordSet", "Value %d (%s) of record %d is not expected type of %s: %#+v", fieldIdx, fieldName, infoIdx, "numeric", info[fieldIdx])
	}

	*field = strconv.Itoa(int(float64Value))
	return ok
}

func assertStringValue(field *string, info []interface{}, fieldIdx int, fieldName string, infoIdx int, logger boshlog.Logger) bool {
	if fieldIdx < 0 {
		return false
	}

	var ok bool
	*field, ok = info[fieldIdx].(string)

	if !ok {
		logger.Warn("RecordSet", "Value %d (%s) of record %d is not expected type of %s: %#+v", fieldIdx, fieldName, infoIdx, "string", info[fieldIdx])
	}

	return ok
}
