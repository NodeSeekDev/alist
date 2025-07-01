package sdk

import "resty.dev/v3"

type RestyOption func(request *resty.Request)

func ReqWithJson(json any) RestyOption {
	return func(request *resty.Request) {
		request.
			SetHeader("Content-Type", "application/json").
			SetBody(json)
	}
}

func ReqWithForm(form Form) RestyOption {
	return func(request *resty.Request) {
		request.
			SetFormData(removeEmptyForm(form))
	}
}

func ReqWithResp(v any) RestyOption {
	return func(request *resty.Request) {
		request.SetResult(v)
	}
}

func ReqWithQuery(query Form) RestyOption {
	return func(request *resty.Request) {
		request.SetQueryParams(removeEmptyForm(query))
	}
}

func ReqWithUA(ua string) RestyOption {
	return func(request *resty.Request) {
		request.SetHeader("User-Agent", ua)
	}
}

func removeEmptyForm(form Form) Form {
	f := Form{}
	for k, v := range form {
		if v == "" {
			continue
		}
		f[k] = v
	}
	return f
}
