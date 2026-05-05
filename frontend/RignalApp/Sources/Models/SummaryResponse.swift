import Foundation

struct SummaryResponse: Codable {
    let date: String
    let conditionScore: Int
    let focusPeakStart: Date?
    let focusPeakEnd: Date?
    let recommendBedtime: Date?
    let sleepDebtMinutes: Int
}
