package label

const (
	CommandList = `Список доступных команд:

*/get C1/C2* - получить Binance-курс для пары C1/C2 (например, */get USDT/BTC*).
	
*/getlist* - получить список поддерживаемых Binance курсов.

*/getall* - получить документ с релевантной информацией о всех курсах.
`
	Start = `Привет. Я умею присылать информацию о релевантных курсах на разных биржах!

` + CommandList
	Menu = `Я могу показать курсы валют с Binance или отправить XLS файл с многими курсами.
`
	GetBinancePairsList = `Выберите интересующую Вас пару из списка или отправьте сообщение в таком же формате.
`
	GetAllBinancePrices = `Документ с последней информацией о курсах с Binance.
`
	GetAllBinancePricesFailed = `Произошла ошибка :( Попробуйте позже.
`
	UnknownCommand = `Неизвестная команда :( Попробуйте позже.
` + CommandList
	GetBinancePrice = `Текущая цена %s: %.8f
`
	GetBinancePriceFailed = `Не удалось получить информацию с Binance. Попробуйте позже.
`
	UnknownPair = `Бот не поддерживает данную пару.
`
	GetList = `Список поддерживаемых Binance курсов:

`
	GetCorrection = `Уточните, курс какой именно пары Вы хотите получить.

Например, */get BTC/USDT*
`
)
