package i18n

import (
	"github.com/gin-gonic/gin/binding"
	locale "github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"

	validator "github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type ErrorMap map[string]map[string]string

var (
	fieldTranslation = make(map[string]ErrorMap)
	_translator      = allocateTranslator([]string{"en", "zh"})
	defaultValidator = binding.Validator.Engine().(*validator.Validate)
)

// RegisterFieldMapping register field mapping
func RegisterFieldMapping(fieldMapping map[string]ErrorMap) {
	if fieldMapping == nil {
		return
	}
	for lng, trans := range fieldMapping {
		if existLngTrans, ok := fieldTranslation[lng]; ok {
			for path, msg := range trans {
				existLngTrans[path] = msg
			}
			fieldTranslation[lng] = existLngTrans
		} else {
			fieldTranslation[lng] = trans
		}
	}
}

// GetFieldMapping return the global private field mapping
func GetFieldMapping() map[string]ErrorMap {
	return fieldTranslation
}

type Translator struct {
	Uni     *ut.UniversalTranslator
	locales []string
}

// NewTranslator return Translator
// Notice locales first element is fallback Translator
func allocateTranslator(locales []string) *Translator {
	if len(locales) == 0 {
		panic("must assign locales")
	}

	trans := make([]locale.Translator, len(locales))
	for index, lng := range locales {
		switch lng {
		case "zh":
			trans[index] = zh.New()
		case "en":
			trans[index] = en.New()
		}
	}

	uni := ut.New(trans[0], trans...)
	t := &Translator{
		Uni:     uni,
		locales: locales,
	}

	return t
}

func GetTranslator() *Translator {
	return _translator
}

func init() {
	defaultValidator.SetTagName("validate")
	//registerDefaultTranslations()
}

// registerDefaultTranslations register default tanslations
func registerDefaultTranslations() {
	for _, locale := range _translator.locales {
		switch locale {
		case "zh":
			zh, _ := _translator.Uni.GetTranslator("zh")
			zh_translations.RegisterDefaultTranslations(defaultValidator, zh)
		case "en":
			en, _ := _translator.Uni.GetTranslator("en")
			en_translations.RegisterDefaultTranslations(defaultValidator, en)
		}
	}
}

// RegisterCustomerTranslation register customer translations
func RegisterCustomerTranslation(tag string, fn validator.Func, templates map[string]string) {
	defaultValidator.RegisterValidation(tag, fn)
	customTag := make(map[string]ut.Translator)
	for _, locale := range _translator.locales {
		switch locale {
		case "zh":
			zh, _ := _translator.Uni.GetTranslator("zh")
			customTag[locale] = zh
		case "en":
			en, _ := _translator.Uni.GetTranslator("en")
			customTag[locale] = en
		}
	}

	for locale, trans := range customTag {
		defaultValidator.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
			return ut.Add(tag, templates[locale], true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
	}
}
