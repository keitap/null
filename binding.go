package null

import (
	"strconv"

	"github.com/mholt/binding"
)

func (v *Bool) Bind(fieldName string, strVals []string, errs binding.Errors) binding.Errors {
	if len(strVals) == 0 {
		v.Valid = false
		return errs
	}

	val, err := strconv.ParseBool(strVals[0])
	if err != nil {
		errs.Add([]string{fieldName}, binding.DeserializationError, err.Error())
	} else {
		v.SetValid(val)
	}
	return errs
}

func (v *Int) Bind(fieldName string, strVals []string, errs binding.Errors) binding.Errors {
	if len(strVals) == 0 {
		v.Valid = false
		return errs
	}

	val, err := strconv.ParseInt(strVals[0], 10, 64)
	if err != nil {
		errs.Add([]string{fieldName}, binding.DeserializationError, err.Error())
	} else {
		v.SetValid(val)
	}
	return errs
}

func (v *Float) Bind(fieldName string, strVals []string, errs binding.Errors) binding.Errors {
	if len(strVals) == 0 {
		v.Valid = false
		return errs
	}

	val, err := strconv.ParseFloat(strVals[0], 64)
	if err != nil {
		errs.Add([]string{fieldName}, binding.DeserializationError, err.Error())
	} else {
		v.SetValid(val)
	}
	return errs
}

func (v *String) Bind(fieldName string, strVals []string, errs binding.Errors) binding.Errors {
	if len(strVals) == 0 {
		v.Valid = false
	} else {
		v.SetValid(strVals[0])
	}
	return errs
}
