import SwiftUI
import WidgetKit

struct WidgetEntryView: View {
    var entry: RignalWidgetProvider.Entry

    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text("Rignal")
                .font(.caption2.bold())
                .foregroundStyle(.secondary)

            Text("\(entry.conditionScore)")
                .font(.system(size: 42, weight: .bold))
                .foregroundStyle(scoreColor(entry.conditionScore))

            Text(scoreLabel(entry.conditionScore))
                .font(.caption.bold())
                .foregroundStyle(.secondary)

            Spacer()

            if let start = entry.focusPeakStart, let end = entry.focusPeakEnd {
                VStack(alignment: .leading, spacing: 2) {
                    Text("集中ピーク")
                        .font(.caption2)
                        .foregroundStyle(.secondary)
                    Text("\(timeString(start))–\(timeString(end))")
                        .font(.caption.bold())
                }
            } else {
                Text("データなし")
                    .font(.caption2)
                    .foregroundStyle(.tertiary)
            }
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
        .padding(12)
    }

    private func scoreColor(_ score: Int) -> Color {
        score >= 80 ? .green : score >= 60 ? .orange : .red
    }

    private func scoreLabel(_ score: Int) -> String {
        score >= 80 ? "Good" : score >= 60 ? "Fair" : "Low"
    }

    private func timeString(_ date: Date) -> String {
        let f = DateFormatter()
        f.dateFormat = "HH:mm"
        f.timeZone = .current
        return f.string(from: date)
    }
}
