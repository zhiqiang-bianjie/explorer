package service

import (
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/orm/document"
	"time"
)

type IService struct {
	BaseService
}

func (service *IService) GetModule() Module {
	return Iservice
}

var svcDefList = []model.SvcDef{
	{
		Hash:              "42C70BE32F4DFBBDD14ECEC395A313582792F7FB10C042D10F2AF6F0B7811374",
		Name:              "BLSQ001",
		ChainId:           "zone1",
		Author:            "iaa1q7602ujxxx0urfw7twm0uk5m7n6l9gqsallwar",
		AuthorDescription: "金融公司1",
		Description:       "保理申请服务定义",
		Status:            "Enable",
		IDLContent: `syntax = "proto2";
option java_package = "com.lhc.protobuf";
option java_outer_classname = "AddressBookProtos";

message Person {
  required string name = 1;
  required int32 id = 2;
  optional string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    required string number = 1;
    optional PhoneType type = 2 [default = HOME];
  }
  repeated PhoneNumber phones = 4;
}

message AddressBook {
  repeated Person people = 1;
}`,
	},
	{
		Hash:              "42C70BE32F4DFBBDD14ECEC395A313582792F7FB10C042D10F2AF6F0B7811374",
		Name:              "CDHPYZ01",
		ChainId:           "zone2",
		Author:            "iaa1as2h0y7x0y6ledn3nv8fyqxnpu28h8rywx2svq",
		AuthorDescription: "金融公司2",
		Description:       "承兑汇票验证服务定义",
		Status:            "Disable",
	},
	{
		Hash:              "7D6661A070EA1008ADC870B3D08B6437C3762F466DC85C8508ACC97D9AAEA1CD",
		Name:              "FPYZ001",
		ChainId:           "zone3",
		Author:            "iaa1mn8mcxvj6hum7vqxy86etatz6prad3d0l8quhd",
		AuthorDescription: "金融公司3",
		Description:       "电子发票验证服务定义",
		Status:            "Disable",
	},
}

var SvcBindList = []model.SvcBind{
	{
		Hash:        "8F4D6FA7CD84C46B6DD7B54445EC95604E21403075E4F776BD283082CD79B2E3",
		DefName:     "BLSQ001",
		DefChainID:  "zone1",
		BindChainID: "zone2",
		Provider:    "iaa1uex02kyx02qap5lrakn5xv0glq49e2ww54d7ka",
		BindingType: "global",
		Prices: document.Coins{{
			Denom:  "iris",
			Amount: 1.5,
		}},
		Level: model.Level{
			AvgRspTime: 2000,
			UsableTime: 9999,
		},
		Available: true,
	},
	{
		Hash:        "4EC915ABB97518DA3AE8BB9A27B51B6EEECEB1CE725CE2024F3BE5F2194083F0",
		DefName:     "BLSQ001",
		DefChainID:  "zone1",
		BindChainID: "zone2",
		Provider:    "iaa1uex02kyx02qap5lrakn5xv0glq49e2ww54d7ka",
		BindingType: "local",
		Prices: document.Coins{{
			Denom:  "iris-atto",
			Amount: 1000000000000000000,
		}, {
			Denom:  "atom",
			Amount: 1,
		}},
		Level: model.Level{
			AvgRspTime: 2500,
			UsableTime: 9900,
		},
		Available: false,
	},
	{
		Hash:        "FE10C0715632F1BED21E47D961822CE33A6A20BFE479546FF6C9D0EAF55721B1",
		DefName:     "CDHPYZ01",
		DefChainID:  "zone1",
		BindChainID: "zone2",
		Provider:    "iaa1q7602ujxxx0urfw7twm0uk5m7n6l9gqsallwar",
		BindingType: "local",
		Prices: document.Coins{{
			Denom:  "iris",
			Amount: 1.5,
		}},
		Available: false,
	},
}

var svcTxList = []model.SvcTx{
	{
		Hash:        "5B1E3051180FD24A878C888F22E58D1B270835978BA8BAEB32F5153F65A98121",
		ReqId:       "35394-35396-1",
		TxType:      "service_call",
		SendAddr:    "iaa163s33re5y7xhuu05r5d8ekwkpdzsjtue73qmpt",
		ReceiveAddr: "iaa1z29lpxul0wtewwmkt9y47r8pxdp77jg2p2w8a8",
		Height:      353944,
		Data:        "******",
		Time:        time.Now(),
	},
	{
		Hash:        "B963F8CD83225B60F39E04D96CC21BDE3F8B04A5CAA5F3688FBF9263707954D9",
		ReqId:       "35394-35396-1",
		TxType:      "service_respond",
		SendAddr:    "iaa1z29lpxul0wtewwmkt9y47r8pxdp77jg2p2w8a8",
		ReceiveAddr: "iaa163s33re5y7xhuu05r5d8ekwkpdzsjtue73qmpt",
		Height:      353945,
		Data:        "******",
		Time:        time.Now(),
	},
}

func (service *IService) QueryList(page, size int) interface{} {
	return svcDefList
}

func (service *IService) Query(name, defChainId string) interface{} {
	var svcInfo IServiceInfo
	for _, svc := range svcDefList {
		if svc.Name == name && svc.ChainId == defChainId {
			svcInfo.SvcDef = svc
		}
	}
	var svcBindList []model.SvcBind
	for _, svcBind := range SvcBindList {
		if svcBind.DefName == name && svcBind.DefChainID == defChainId {
			svcBindList = append(svcBindList, svcBind)
		}
	}
	svcInfo.SvcBindList = svcBindList
	svcInfo.SvcTxList = svcTxList
	return svcInfo
}

type IServiceInfo struct {
	model.SvcDef
	SvcBindList []model.SvcBind `json:"svc_bind_list"`
	SvcTxList   []model.SvcTx   `json:"svc_tx_list"`
}
