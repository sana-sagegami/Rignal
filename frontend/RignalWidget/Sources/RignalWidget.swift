import WidgetKit
import SwiftUI

struct WidgetEntry: TimelineEntry {
    let date: Date
    let conditionScore: Int
    let focusPeakStart: Date?
    let focusPeakEnd: Date?
}

struct RignalWidgetProvider: TimelineProvider {
    private let appGroupID = "group.dev.tuist.RignalApp"

    func placeholder(in context: Context) -> WidgetEntry {
        WidgetEntry(date: .now, conditionScore: 75, focusPeakStart: nil, focusPeakEnd: nil)
    }

    func getSnapshot(in context: Context, completion: @escaping (WidgetEntry) -> Void) {
        completion(currentEntry())
    }

    func getTimeline(in context: Context, completion: @escaping (Timeline<WidgetEntry>) -> Void) {
        let entry = currentEntry()
        // 翌朝 7:00 に再取得
        var components = Calendar.current.dateComponents([.year, .month, .day], from: .now)
        components.day! += 1
        components.hour = 7
        components.minute = 0
        let refreshDate = Calendar.current.date(from: components) ?? Date().addingTimeInterval(3600 * 12)
        completion(Timeline(entries: [entry], policy: .after(refreshDate)))
    }

    private func currentEntry() -> WidgetEntry {
        let defaults = UserDefaults(suiteName: appGroupID)
        let score = defaults?.integer(forKey: "conditionScore") ?? 0
        let startInterval = defaults?.double(forKey: "focusPeakStart") ?? 0
        let endInterval = defaults?.double(forKey: "focusPeakEnd") ?? 0
        return WidgetEntry(
            date: .now,
            conditionScore: score,
            focusPeakStart: startInterval > 0 ? Date(timeIntervalSince1970: startInterval) : nil,
            focusPeakEnd: endInterval > 0 ? Date(timeIntervalSince1970: endInterval) : nil
        )
    }
}

struct RignalWidget: Widget {
    let kind = "RignalWidget"

    var body: some WidgetConfiguration {
        StaticConfiguration(kind: kind, provider: RignalWidgetProvider()) { entry in
            WidgetEntryView(entry: entry)
                .containerBackground(.fill.tertiary, for: .widget)
        }
        .configurationDisplayName("Rignal")
        .description("コンディションスコアと集中ピーク時間帯を表示します")
        .supportedFamilies([.systemSmall])
    }
}
