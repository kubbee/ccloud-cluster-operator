package services

/*
 * This function get the output console and converts to string
 */
func getOutput(outs []byte) (string, bool) {
	if len(outs) > 0 {
		return string(outs), true
	}
	return "", false
}
