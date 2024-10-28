package demo

import "time"

// AddMonthsPreservingEndOfMonth は、指定された月数を日付 t に加算します。元の日付が月の末日付近であった場合、は日付の調整を行います。
// 標準の time.AddDate 関数では、月末のない月の処理が行えないためこの関数を作成しました。
// 例えば、1月31日に1ヶ月をtime.AddDateで加算すると、3月3日となりますが、
// この関数では加算後の月の末日（2月の最終日28日または29日）に調整します
//
// 引数:
// - t: 基準の日付
// - months: 加算する月数
//
// 返り値:
// - time.Time: 月を加算した後の新しい日付
func AddMonthsPreservingEndOfMonth(t time.Time, months int) time.Time {
	// 基準日の日付部分を取得
	day := t.Day()

	// 加算後の日付の最終日を取得
	year, month, dayLimit := t.AddDate(0, months+1, -day).Date()

	// 基準日が月末付近だった際の日付の調整
	// 左辺は基準日の日付が30日や31日で加算後の月に30日や31日が存在しない場合に対応するため
	// 右辺は基準日の月の月末が30日で加算後の月の月末が31日の場合に対応するため
	if day > dayLimit || t.Month() != t.AddDate(0, 0, 1).Month() {
		day = dayLimit
	}

	h, m, s := t.Clock()
	return time.Date(year, month, day, h, m, s, t.Nanosecond(), t.Location())
}
