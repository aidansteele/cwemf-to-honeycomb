package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransformer(t *testing.T) {
	tt := &transformer{}

	input := events.KinesisFirehoseEvent{}
	err := json.Unmarshal(goodPayload, &input)
	require.NoError(t, err)

	resp, err := tt.handle(context.Background(), &input)
	require.NoError(t, err)

	j, _ := json.Marshal(resp)
	fmt.Println(string(j))
}

var goodPayload = []byte(`
{
  "invocationId": "7a066bf5-d337-452c-8e74-a0ba185c51e3",
  "deliveryStreamArn": "arn:aws:firehose:ap-southeast-2:607481581596:deliverystream/cwemf-to-honeycomb-Firehose-mQpnsQxMhQW3",
  "sourceKinesisStreamArn": "",
  "region": "ap-southeast-2",
  "records": [
    {
      "recordId": "49624128254508124554132954738458680784257569011217203202000000",
      "approximateArrivalTimestamp": 1637537167704,
      "data": "H4sIAAAAAAAAAL2TbYvaQBSF/8t8Vpw7L3dm/CZdXUrZUlBaignLmEzd0Lw1iVoR/3tvIrJtdwv7pUII4dw595w8JCdWhLb127A61oFN2d1sNXt8mC+Xs/s5G7HqUIaGZORGWdB0OSQ5r7b3TbWraTLxh3aS+2KT+klyCMW3cVeNn6oyHJOq2IyXvqjzMP756cO20cU++7H4evEvuyb4ghYILmACMKH7WsTcmjSxxistUNlgN6AMF+iSEDRPbR/e7jZt0mR1l1XlIsu70LRsumYvw1k8JM33oez6IyeWpRQoUYMVzjqwIMBKjg5QIVq0QgonUQipEdEIZZUwSjutJFeSoruMaHX0SmwKKI2Wpmfi3OhKkdafIrbIQp5GbBqxzz7fhYiNIva+rHcdaTT+Ho4wTPf9FIYxaeJZE1dNPmv0eCbxkXhf1qyuZaK/20TsXV7t0i++S54eQtdkSW9Zk+ejJ0/tkzDs3VaEbIi6o11lS0CHg+s4Ju03Z3yOz+eopAJ/MnSca3RoteEWhUGlJElAsKzlIKU2EhVHUEqhJMLwOkMUIPAWDNVbGF7a3JChBeus5lIJh1prK8EoI8FJIzV3hBe0FEYaZ6D/Af/BUNNXeguG+k0MqY34jwzj8y/xFD/ruQQAAA==",
      "kinesisRecordMetadata": {
        "shardId": "",
        "partitionKey": "",
        "sequenceNumber": "",
        "subsequenceNumber": 0,
        "approximateArrivalTimestamp": -6795364578871
      }
    },
    {
      "recordId": "49624128254508124554132954739605951387071852853927608322000000",
      "approximateArrivalTimestamp": 1637537178183,
      "data": "H4sIAAAAAAAAAL2QW4saQRCF/0s/K3b1rap9k6wuIWwIKAnBGZZxpuMOmVtmRo0s/veUI7Jhk8C+RGi64VTVOdXfsyhD1yXbsDo2QUzF3Ww1e3yYL5ez+7kYifpQhZZlJ9EQWD7esVzU2/u23jVcmSSHblIk5SZLJukhlN/GfT1+qqtwTOtyM14mZVOE8c9PH7atLff5j8XXy/yyb0NSsoGSCiYAE77XKpaEWUqYGKucoUAbMCiV82kIVmZ0Du92my5t86bP62qRF31oOzFdiz/DRTwkzfeh6s8tzyLPOFA7C6S8B6uMJelAaeTH8A8lKqutJECjjCOlldTgPSIReI7uc6bV85fEFJxGqxEcocTRlSLbP0dikYcii8Q0Ep+TYhciMYrE+6rZ9axx+Xs4wlDdn6swlFlTL5q6avpFc5E4sfjIvC82q+sy0ettIvGuqHfZl6RPnx5C3+bpeWTNMx8TnmmSNAy+25qRDVF37FV1DHRoXMcxa79Nxqf4dIoqXuAVQwUawEnGw4ykR+fRKqWlc0BoyUrtnPHeOO6X2ti/M0Qw6iYM8S0ML9vcjqE+g3FGMjuUFrUy0pP1QNZZT9oQIXknwWtAax3+g6E1SLdgSG9iOGzz/xjGp19yowp+uQQAAA==",
      "kinesisRecordMetadata": {
        "shardId": "",
        "partitionKey": "",
        "sequenceNumber": "",
        "subsequenceNumber": 0,
        "approximateArrivalTimestamp": -6795364578871
      }
    }
  ]
}
`)