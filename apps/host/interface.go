package host

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func (req *QueryHostRequest) Offset() int64 {
	return (req.PageNumber - 1) * req.PageSize
}

func NewDescribeHostRequestWithID(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		Id: id,
	}
}

func NewPatchUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode: UpdateMode_PATCH,
		Resource:   &Resource{},
		Describe:   &Describe{},
	}
}

func NewPutUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode: UpdateMode_PUT,
		Resource:   &Resource{},
		Describe:   &Describe{},
	}
}
