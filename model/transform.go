package model

import "google.golang.org/protobuf/proto"

func CloneMetadata(m *Metadata) *Metadata {
	o := proto.Clone(m)
	o2 := o.(*Metadata)
	return o2
}

// GenerateOutputEvent generates output event
func GenerateOutputEvent(inputMetadata *Metadata,
	contractId string,
	status string,
	statusCode uint32,
	typeUrl string,
	data []byte) *Event {

	oE := &Event{TypeUrl: typeUrl, Value: data}

	if inputMetadata != nil {
		oE.Metadata = CloneMetadata(inputMetadata)
		if oE.Metadata != nil && oE.Metadata.Headers != nil {
			oE.Metadata.Headers["Content-Type"] = typeUrl
		}

		if oE.Metadata != nil {
			oE.Metadata.Status = status
			oE.Metadata.StatusCode = statusCode

			if oE.Metadata.ContractIdList == nil {
				oE.Metadata.ContractIdList = make([]string, 0)
			}
			oE.Metadata.ContractIdList = append(oE.Metadata.ContractIdList, contractId)
		}
	}

	return oE
}
