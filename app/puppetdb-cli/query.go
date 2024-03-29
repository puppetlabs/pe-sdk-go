package app

import (
	"context"
	"encoding/json"
	"fmt"

	puppetdbjson "github.com/puppetlabs/pe-sdk-go/json"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/log"

	httptransport "github.com/go-openapi/runtime/client"
)

// QueryWithErrorDetails performs a pdb query and prints error details if any
func (puppetDb *PuppetDb) QueryWithErrorDetails(args string) (interface{}, error) {

	resp, err := puppetDb.query(args)
	if err != nil {
		if du, ok := err.(*operations.GetQueryInternalServerError); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[GET /pdb/query/v4][500] getQuery ServerError: %+v", puppetdbjson.PrettyPrintPayload(du.Payload))
			}
		}
		if du, ok := err.(*operations.GetQueryDefault); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[GET /pdb/query/v4][%d] getQuery default: %+v\n%v", du.Code(), du.Payload.Msg, puppetdbjson.PrettyPrintPayload(du.Payload.Details))
			}
		}
		return "", err
	}
	return resp.Payload, err
}

func (puppetDb *PuppetDb) query(args string) (*operations.GetQueryOK, error) {
	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return nil, err
	}

	args, err = convertAST(args)
	if err != nil {
		log.Debug("Cannot convert query to AST")
	}
	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", puppetDb.Token)

	queryParameters := operations.NewGetQueryParamsWithContext(context.Background())
	queryParameters.SetQuery(&args)

	return client.Operations.GetQuery(queryParameters, apiKeyHeaderAuth)
}

func convertAST(args string) (string, error) {

	var data []interface{}
	bytes := []byte(args)
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Debug(fmt.Sprintf("Cannot unmarshal %v to json array", args))
		return args, err
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Debug(fmt.Sprintf("Cannot marshal %v to string", data))
		return args, err
	}
	return string(jsonBytes), nil
}
