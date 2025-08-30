package validation

func HostSchema() Schema {
	return Schema{
		Validators: map[string]Validator{
			"name":  validateName,
			"user":  validateUser,
			"host":  validateHost,
			"port":  validatePort,
			"force": validateForce,
		},
	}
}
