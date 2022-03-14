package ipapi

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ClientFree(t *testing.T) {

	client := New("")
	assert.Equal(t, client.buildURL(""), "http://ip-api.com/json/?fields=status%2Cmessage%2Ccountry%2CcountryCode%2Cregion%2CregionName%2Ccity%2Czip%2Clat%2Clon%2Ctimezone%2Cisp%2Corg%2Cas%2Cquery")
	assert.Equal(t, client.buildURL("1.2.3.4"), "http://ip-api.com/json/1.2.3.4?fields=status%2Cmessage%2Ccountry%2CcountryCode%2Cregion%2CregionName%2Ccity%2Czip%2Clat%2Clon%2Ctimezone%2Cisp%2Corg%2Cas%2Cquery")
	q, err := client.Query(context.Background(), "127.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, q.Status, "fail")
	assert.Equal(t, q.Message, "reserved range")

	q, err = client.Query(context.Background(), "")
	assert.Nil(t, err)
	assert.Equal(t, q.Status, "success")
	assert.Equal(t, q.Message, "")

}

func Test_ClientPro(t *testing.T) {

	client := New("12345")
	assert.Equal(t, client.buildURL(""), "https://pro.ip-api.com/json/?apiKey=12345&fields=status%2Cmessage%2Ccountry%2CcountryCode%2Cregion%2CregionName%2Ccity%2Czip%2Clat%2Clon%2Ctimezone%2Cisp%2Corg%2Cas%2Cquery")
	assert.Equal(t, client.buildURL("1.2.3.4"), "https://pro.ip-api.com/json/1.2.3.4?apiKey=12345&fields=status%2Cmessage%2Ccountry%2CcountryCode%2Cregion%2CregionName%2Ccity%2Czip%2Clat%2Clon%2Ctimezone%2Cisp%2Corg%2Cas%2Cquery")

}

func Test_Client_Funcs(t *testing.T) {

	client := New("")
	assert.True(t, client.fieldAllowed("status"))
	assert.False(t, client.fieldAllowed("nonAllowed"))

	assert.Equal(t, client.fields, "status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query")
	client.SetFields([]string{"status", "message", "query", "invalid1", "invalid2"}, false)
	assert.Equal(t, client.fields, "status,message,query")

	client.SetTimeout(10 * time.Second)
	assert.Equal(t, client.timeout, 10*time.Second)

	client.SetFields([]string{"status", "query", "query", "query", "message"}, true)
	assert.Equal(t, client.fields, "57344")

}
