package utils


type ListAndCount struct {
	Items interface{} `json:"items"`
	Count uint64      `json:"total"`
	Info  interface{} `json:"info"`
}

func Lists(items interface{}, count uint64) *ListAndCount {
	if items == nil {
		return &ListAndCount{
			Items: []struct{}{},
			Count: count,
		}
	}
	return &ListAndCount{
		Items: items,
		Count: count,
	}
	
}
