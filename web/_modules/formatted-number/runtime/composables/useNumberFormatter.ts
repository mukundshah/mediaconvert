import { Decimal } from 'decimal.js'
import { computed } from 'vue'

// type FormatType = 'decimal' | 'currency' | 'unit' | 'percent'
// type UnitDisplay = Intl.NumberFormatOptions['unitDisplay']
// type CurrencyDisplay = Intl.NumberFormatOptions['currencyDisplay'] | 'never'
// type SignDisplay = Intl.NumberFormatOptions['signDisplay'] | 'accounting'

// Global cache for Intl.NumberFormat instances
const formatterCache = new Map<string, Intl.NumberFormat>()

const getFormatter = (locale: string, options: Intl.NumberFormatOptions): Intl.NumberFormat => {
  const cacheKey = `${locale}-${JSON.stringify(options)}`

  if (!formatterCache.has(cacheKey)) {
    formatterCache.set(cacheKey, new Intl.NumberFormat(locale, options))
  }

  return formatterCache.get(cacheKey)!
}

export interface NumberFormatterOptions {
  locale?: string
  options?: Intl.NumberFormatOptions
  nullValue?: string
  zeroValue?: string
}

export function useNumberFormatter(
  options: NumberFormatterOptions = {},
) {
  const {
    locale = 'en-US',
    zeroValue,
    nullValue = '-',
    options: numberFormatOptions = {},
  } = options

  const formatter = getFormatter(locale, numberFormatOptions)

  const format = (value: number | string | Decimal | null | undefined) => {
    const currentValue = value

    if (currentValue === null || currentValue === undefined) {
      return nullValue
    }

    const _value = new Decimal(currentValue)

    if (_value.isZero() && zeroValue) {
      return zeroValue
    }

    return formatter.format(_value.toNumber())

    // if (type === 'currency') {
    //   formatOptions.currency = currency
    //   // formatOptions.currencySign = 'accounting'
    //   // formatOptions.signDisplay = 'never'

    //   if (currencyDisplay === 'never') {
    //     formatOptions.currencyDisplay = 'code'
    //   } else {
    //     formatOptions.currencyDisplay = currencyDisplay
    //   }

    //   if (formatOptions.currency === 'NPR' || formatOptions.currency === 'INR') {
    //     _locale = 'en-IN' // only en-IN has lakhs and crores separators
    //   }
    // } else if (type === 'unit') {
    //   formatOptions.unit = unit
    //   formatOptions.unitDisplay = unitDisplay
    // }

    // if (currencyDisplay === 'never') {
    //   formattedValue = formattedValue.replace(currency || '', '').trim()
    // }

    // // Add brackets for negative values for decimal and unit types
    // if (_value.isNegative()) {
    //   return `(${formattedValue})`
    // }
  }

  return { format }
}
