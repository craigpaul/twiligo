package twiligo

// ConvertDirectionToMessageDirection ...
func ConvertDirectionToMessageDirection(direction string) MessageDirection {
	return map[string]MessageDirection{
		"inbound":        Inbound,
		"outbound-api":   OutboundAPI,
		"outbound-call":  OutboundCall,
		"outbound-reply": OutboundReply,
	}[direction]
}

// ConvertStatusToMessageStatus ...
func ConvertStatusToMessageStatus(status string) MessageStatus {
	return map[string]MessageStatus{
		"accepted":    Accepted,
		"queued":      Queued,
		"sending":     Sending,
		"sent":        Sent,
		"failed":      Failed,
		"delivered":   Delivered,
		"undelivered": Undelivered,
		"receiving":   Receiving,
		"received":    Received,
	}[status]
}

// GetPhoneNumberType will convert a given integer into a PhoneNumberType by the given country. Certain countries do not support all PhoneNumberType values, so this function can be used as a safe mapping based on the values from this document https://support.twilio.com/hc/en-us/articles/223183068-Twilio-international-phone-number-availability-and-their-capabilities. Note: Not all cases are currently supported, but can be amended as necessary.
func GetPhoneNumberType(number int, country string) PhoneNumberType {
	numberType := PhoneNumberType(number)

	if (country == "CA" || country == "US") && numberType == Mobile {
		numberType = Local
	}

	return numberType
}
