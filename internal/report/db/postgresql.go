package db

import (
	"context"
	"strings"
	"wb-report-downloader/internal/report"
	"wb-report-downloader/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, report *report.ReportDetailes) error {
	q := `
	INSERT INTO wb_reports_details_v2 ("№", "Номер поставки", "Предмет", "Код номенклатуры", "Бренд", "Артикул поставщика",
		"Название", "Размер", "Баркод", "Тип документа", "Обоснование для оплаты",
		"Дата заказа покупателем", "Дата продажи", "Кол-во", "Цена розничная",
		"Вайлдберриз реализовал Товар (Пр)", "Согласованный продуктовый дискон", "Промокод %",
		"Итоговая согласованная скидка", "Цена розничная с учетом согласова",
		"Скидка постоянного Покупателя (СП", "Размер кВВ, %", "Размер  кВВ без НДС, % Базовый",
		"Итоговый кВВ без НДС, %", "Вознаграждение с продаж до вычета ",
		"Вознаграждение Вайлдберриз (ВВ), б",
		"НДС с Вознаграждения Вайлдберриз", "К перечислению Продавцу за реализ",
		"Количество доставок", "Количество возврата", "Услуги по доставке товара покупат",
		"Штрафы", "Доплаты", "Обоснование штрафов и доплат", "Стикер МП", "Номер офиса",
		"Наименование офиса доставки", "ИНН партнера", "Партнер", "Склад", "Страна",
		"Тип коробов", "Номер таможенной декларации", "ШК", "Rid", "Srid",
		"Возмещение за выдачу и возврат тов", "Возмещение расходов по эквайрингу", "Наименование банка эквайринга",
		report_id, seller_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
		$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, 
		$40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51)
	`
	// fmt.Printf("SQL Query: %s", formatQuery(q))

	_, err := r.client.Exec(ctx, q, 
		report.Number,
		report.SupplyNumber,
		report.Subject,
		report.ItemCode,
		report.Brand,
		report.SuppliersArticle,
		report.Name,
		report.Size,
		report.Barcode,
		report.DocumentType,
		report.JustificationForPayment,
		report.DateOfTheOrderByTheBuyer,
		report.DateOfSale,
		report.Count,
		report.RetailPrice,
		report.WildberriesSoldTheProductPr,
		report.AgreedGroceryDiscount,
		report.Promocode,
		report.FinalAgreedDiscount,
		report.RetailPriceIncludingTheAgreedDiscount,
		report.RegularCustomerDiscountSPP,
		report.KVVSizePercent,
		report.KVVSizeWithoutVATBasicPercent,
		report.FinalKVVWithoutVATPercent,
		report.RemunerationFromSalesBeforeDeductionDfAttorneysServicesWithoutVAT,
		// report.ReimbursementOfAttorneysExpenses,
		report.WildberriesRemunerationBBWithoutVAT,
		report.VATOnWildberriesRemuneration,
		report.ToTransferToTheSellerForTheSoldGoods,
		report.NumberOfDeliveries,
		report.RefundAmount,
		report.ServicesForTheDeliveryOfGoodsToTheBuyer,
		report.Fines,
		report.Surcharges,
		report.JustificationOfFinesAndSurcharges,
		report.StickerMP,
		report.OfficeNumber,
		report.NameOfTheDeliveryOffice,
		report.PartnersINN,
		report.Partner,
		report.Warehouse,
		report.Country,
		report.TypeOfBoxes,
		report.CustomsDeclarationNumber,
		report.ShK,
		report.Rid,
		report.Srid,

		report.RefundForTheDeliveryAndReturnOfGoodsToThePVZ,
		report.ReimbursementOfAcquiringCosts,
		report.NameOfAcquiringBank,

		report.ReportID,
		report.SellerID)
	
	return err
}

func NewRepository(client postgresql.Client) report.Repository {
	return &repository{
		client: client,
	}
}
