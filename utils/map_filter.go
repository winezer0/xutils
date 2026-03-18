package utils

// FilterMapByKeys  过滤掉那些键已经存在于 filteredKeys 列表中的条
func FilterMapByKeys(sourceMap map[string]map[string]string, filteredKeys []string) map[string]map[string]string {
	if len(filteredKeys) == 0 {
		return sourceMap
	}

	set := make(map[string]struct{}, len(filteredKeys))
	for _, k := range filteredKeys {
		set[k] = struct{}{}
	}
	out := make(map[string]map[string]string)
	for k, v := range sourceMap {
		if _, ok := set[k]; !ok {
			out[k] = v
		}
	}
	return out
}
