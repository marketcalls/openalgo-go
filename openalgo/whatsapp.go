package openalgo

// WhatsApp sends a WhatsApp message via the OpenAlgo paired device.
//
// One call handles every common case - send to self, to a single number, to
// a small broadcast (max 5 recipients), with text, image, or document
// attachments. The OpenAlgo server must already be paired to a WhatsApp
// account from the "/whatsapp" admin page; pairing is not exposed via the
// API.
//
// message is the plain text body (max 4096 characters). Pass "" if you are
// only sending an image/document with a caption.
//
// optionalParams (pass a single map, or omit) may include:
//   - "to" (string or []string): recipient(s). A single E.164 digit string
//     (e.g. "919876543210") for one recipient, or a []string of up to 5 for a
//     small broadcast (anything beyond 5 is dropped server-side).
//   - "username" (string): OpenAlgo username - resolved via the linked-users
//     table on the server. Takes precedence over "to" if both are set.
//   - "image" (string): server-local path to an image file.
//   - "document" (string): server-local path to a document (PDF, CSV, ...).
//   - "caption" (string): caption for the image, or follow-up text for the document.
//   - "filename" (string): override the document's display name on the recipient's device.
//   - "wait_for_delivery" (bool): default true. When true, blocks until WhatsApp
//     confirms delivery and returns a per-recipient report. Set to false for
//     fire-and-forget; the response is a generic "queued" acknowledgement.
//
// If neither "to" nor "username" is supplied, the message defaults to
// self=true (sent to the paired device's own number).
//
// Example:
//
//	// Send to self
//	client.WhatsApp("Build #482 deployed. P&L: +1.2%")
//
//	// Send to a single number
//	client.WhatsApp("Are you free for a quick call?", map[string]interface{}{
//	    "to": "919876543210",
//	})
//
//	// Small broadcast with an image
//	client.WhatsApp("NIFTY end-of-day chart", map[string]interface{}{
//	    "to":    []string{"919876543210", "919812345678"},
//	    "image": "/srv/charts/nifty_eod.png",
//	})
func (c *Client) WhatsApp(message string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}

	waitForDelivery := true
	recipientSet := false

	if len(optionalParams) > 0 {
		params := optionalParams[0]

		if username, ok := params["username"].(string); ok && username != "" {
			payload["username"] = username
			recipientSet = true
		} else if to, ok := params["to"]; ok {
			switch v := to.(type) {
			case []string:
				if len(v) > 0 {
					payload["phones"] = v
					recipientSet = true
				}
			case string:
				if v != "" {
					payload["phone"] = v
					recipientSet = true
				}
			}
		}

		if image, ok := params["image"].(string); ok && image != "" {
			payload["image_path"] = image
		}
		if document, ok := params["document"].(string); ok && document != "" {
			payload["document_path"] = document
		}
		if caption, ok := params["caption"].(string); ok && caption != "" {
			payload["caption"] = caption
		}
		if filename, ok := params["filename"].(string); ok && filename != "" {
			payload["filename"] = filename
		}
		if wait, ok := params["wait_for_delivery"].(bool); ok {
			waitForDelivery = wait
		}
	}

	if !recipientSet {
		payload["self"] = true
	}

	if message != "" {
		payload["message"] = message
	}
	payload["wait_for_delivery"] = waitForDelivery

	return c.makeRequest("POST", "whatsapp/notify", payload)
}
