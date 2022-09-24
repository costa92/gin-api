package options

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServingOptions.Validate()...)
	errs = append(errs, o.MysqlConfig.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Jwt.Validate()...)
	errs = append(errs, o.Logger.Validate()...)
	return errs
}
