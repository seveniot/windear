package ex

import "errors"

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/5
 * 
*/

var (
	EmptyPacket = errors.New("the packet is empty")
	IncorrectRemainLength = errors.New("the remain length is incorrect")
	IncorrectMSBValue  = errors.New("the msb value is incorrect")
	IncorrectLSBValue   = errors.New("the lsb value is incorrect")
	MqttIdentifierNotFound  = errors.New("the MQTT identifier is not found in the connect packet")
	UnsupportedMqttVersion = errors.New("the MQTT version is not supported")
	IllegalPayloadData  = errors.New("the payload data is illegal")
	IllegalQOSLEVEL = errors.New("the qos level is illegal")
	LackOfMessageIdentifier = errors.New("the message need an identifier")
	IncorrectSubscribeTopic  = errors.New("the subscribe topic is incorrect")
	MalformedPacketType = errors.New("the packet type is malformed")
	MalformedPacketFlag = errors.New("the packet flag is malformed")
	IncorrectRemainLengthMultiplier = errors.New("the remain length multiplier is incorrect")

	IncorrectSessionlId = errors.New("the sessionId is incorrect")

	ConnNoAuthFound = errors.New("no auth information found in connect packet")
	ConnAuthFail = errors.New("auth fail during connect")

	ClientDisconnect = errors.New("client disconnected")

	AgentReturnFail = errors.New("agent return fail")
)