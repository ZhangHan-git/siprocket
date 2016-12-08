package siprocket

/*
 RFC 3261 - https://www.ietf.org/rfc/rfc3261.txt

INVITE sip:01798300765@87.252.61.202;user=phone SIP/2.0
SIP/2.0 200 OK

*/

type sipReq struct {
	Method     []byte // Sip Method eg INVITE etc
	UriType    string // Type of URI sip, sips, tel etc
	StatusCode []byte // Status Code
	User       []byte // User part
	Host       []byte // Host part
	Port       []byte // Port number
	UserType   []byte // User Type
	Src        []byte // Full source if needed
}

func parseSipReq(v []byte, out *sipReq) {

	pos := 0
	state := 0

	// Init the output area
	out.UriType = ""
	out.Method = nil
	out.StatusCode = nil
	out.User = nil
	out.Host = nil
	out.Port = nil
	out.UserType = nil
	out.Src = nil

	// Keep the source line if needed
	if keep_src {
		out.Src = v
	}

	// Loop through the bytes making up the line
	for pos < len(v) {
		// FSM
		switch state {
		case FIELD_NULL:
			if v[pos] >= 'A' && v[pos] <= 'S' && pos == 0 {
				state = FIELD_METHOD
				continue
			}

		case FIELD_METHOD:
			if v[pos] == ' ' || pos > 9 {
				if string(out.Method) == "SIP/2.0" {
					state = FIELD_STATUS
				} else {
					state = FIELD_BASE
				}
				pos++
				continue
			}
			out.Method = append(out.Method, v[pos])

		case FIELD_BASE:
			if v[pos] != ' ' {
				// Not a space so check for uri types
				if string(v[pos:pos+4]) == "sip:" {
					state = FIELD_USER
					pos = pos + 4
					out.UriType = "sip"
					continue
				}
				if string(v[pos:pos+5]) == "sips:" {
					state = FIELD_USER
					pos = pos + 5
					out.UriType = "sips"
					continue
				}
				if string(v[pos:pos+4]) == "tel:" {
					state = FIELD_USER
					pos = pos + 4
					out.UriType = "tel"
					continue
				}
				// Look for a Tag identifier
				// if string(v[pos:pos+4]) == "tag=" {
				//     state = FIELD_TAG
				//     pos = pos + 4
				//     continue
				// }
				// Look for a User Type identifier

				if string(v[pos:pos+5]) == "user=" {
					state = FIELD_USERTYPE
					pos = pos + 5
					continue
				}
			}
		case FIELD_USER:
			if v[pos] == ':' {
				state = FIELD_PORT
				pos++
				continue
			}
			if v[pos] == ';' || v[pos] == '>' {
				state = FIELD_BASE
				pos++
				continue
			}
			if v[pos] == '@' {
				state = FIELD_HOST
				out.User = out.Host // Move host to user
				pos++
				continue
			}
			out.Host = append(out.Host, v[pos]) // Append to host for now

		case FIELD_HOST:
			if v[pos] == ':' {
				state = FIELD_PORT
				pos++
				continue
			}
			if v[pos] == ';' || v[pos] == '>' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.Host = append(out.Host, v[pos])

		case FIELD_PORT:
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.Port = append(out.Port, v[pos])

		case FIELD_USERTYPE:
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.UserType = append(out.UserType, v[pos])

		case FIELD_STATUS:
			if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
				state = FIELD_BASE
				pos++
				continue
			}
			out.StatusCode = append(out.StatusCode, v[pos])

			// case FIELD_TAG:
			// 	if v[pos] == ';' || v[pos] == '>' || v[pos] == ' ' {
			// 		state = FIELD_BASE
			// 		pos++
			// 		continue
			// 	}
			// 	out.Tag = append(out.Tag, v[pos])

		}
		pos++
	}
}
