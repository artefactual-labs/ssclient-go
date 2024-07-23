package models
type LocationPurpose int

const (
    AR_LOCATIONPURPOSE LocationPurpose = iota
    AS_LOCATIONPURPOSE
    CP_LOCATIONPURPOSE
    DS_LOCATIONPURPOSE
    SD_LOCATIONPURPOSE
    SS_LOCATIONPURPOSE
    BL_LOCATIONPURPOSE
    TS_LOCATIONPURPOSE
    RP_LOCATIONPURPOSE
)

func (i LocationPurpose) String() string {
    return []string{"AR", "AS", "CP", "DS", "SD", "SS", "BL", "TS", "RP"}[i]
}
func ParseLocationPurpose(v string) (any, error) {
    result := AR_LOCATIONPURPOSE
    switch v {
        case "AR":
            result = AR_LOCATIONPURPOSE
        case "AS":
            result = AS_LOCATIONPURPOSE
        case "CP":
            result = CP_LOCATIONPURPOSE
        case "DS":
            result = DS_LOCATIONPURPOSE
        case "SD":
            result = SD_LOCATIONPURPOSE
        case "SS":
            result = SS_LOCATIONPURPOSE
        case "BL":
            result = BL_LOCATIONPURPOSE
        case "TS":
            result = TS_LOCATIONPURPOSE
        case "RP":
            result = RP_LOCATIONPURPOSE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLocationPurpose(values []LocationPurpose) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LocationPurpose) isMultiValue() bool {
    return false
}
