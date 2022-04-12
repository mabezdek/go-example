package dataProvider

type DataProvider struct {
	Data map[string]interface{}
}

func (d *DataProvider) Set(k string, v interface{}) {
	if d.Data == nil {
		d.Data = make(map[string]interface{})
	}

	d.Data[k] = v
}

func (d *DataProvider) Get(k string) interface{} {
	return d.Data[k]
}

func (d *DataProvider) Decode(k string, result interface{}) error {
	return Decode(d.Get(k), result)
}
