import SwiftUI

struct DashboardView: View {
    var authManager: GoogleAuthManager
    @State private var viewModel = DashboardViewModel()

    var body: some View {
        NavigationStack {
            if !authManager.isSignedIn {
                signInView
                    .navigationTitle("Rignal")
            } else {
                Group {
                    if viewModel.isLoading {
                        ProgressView("読み込み中...")
                    } else if let summary = viewModel.summary {
                        summaryContent(summary)
                    } else if let error = viewModel.errorMessage {
                        VStack(spacing: 16) {
                            Text("エラー: \(error)")
                                .foregroundStyle(.red)
                                .multilineTextAlignment(.center)
                                .padding()
                            Button("再読み込み") {
                                Task { await viewModel.load(authManager: authManager) }
                            }
                            .buttonStyle(.bordered)
                        }
                    } else {
                        VStack(spacing: 16) {
                            Text("データなし")
                                .foregroundStyle(.secondary)
                            Button("再読み込み") {
                                Task { await viewModel.load(authManager: authManager) }
                            }
                            .buttonStyle(.bordered)
                        }
                    }
                }
                .navigationTitle("Rignal")
                .toolbar {
                    ToolbarItem(placement: .topBarTrailing) {
                        Button("サインアウト") { authManager.signOut() }
                            .font(.footnote)
                    }
                    #if DEBUG
                    ToolbarItem(placement: .topBarLeading) {
                        Button {
                            Task { await viewModel.triggerAnalysis(authManager: authManager) }
                        } label: {
                            Label("分析実行", systemImage: "arrow.clockwise.circle")
                                .font(.footnote)
                        }
                    }
                    #endif
                }
                .task { await viewModel.load(authManager: authManager) }
                .refreshable { await viewModel.load(authManager: authManager) }
            }
        }
    }

    private var signInView: some View {
        VStack(spacing: 24) {
            Spacer()
            Text("Rignal")
                .font(.largeTitle.bold())
            Text("Google アカウントでサインインして\nカレンダーを連携します")
                .multilineTextAlignment(.center)
                .foregroundStyle(.secondary)
            Button {
                Task { try? await authManager.signIn() }
            } label: {
                Label("Google でサインイン", systemImage: "person.crop.circle.badge.checkmark")
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(.borderedProminent)
            .padding(.horizontal, 40)
            Spacer()
        }
        .padding()
    }

    @ViewBuilder
    private func summaryContent(_ summary: SummaryResponse) -> some View {
        ScrollView {
            VStack(spacing: 20) {
                conditionCard(score: summary.conditionScore)
                if let title = viewModel.nextEventTitle, let time = viewModel.nextEventTime {
                    nextEventCard(title: title, time: time)
                }
                focusPeakCard(start: summary.focusPeakStart, end: summary.focusPeakEnd)
                if let bedtime = summary.recommendBedtime {
                    bedtimeCard(bedtime: bedtime)
                }
                sleepDebtCard(minutes: summary.sleepDebtMinutes)
            }
            .padding()
        }
    }

    private func conditionCard(score: Int) -> some View {
        VStack(spacing: 8) {
            Text("コンディションスコア")
                .font(.headline)
            Text("\(score)")
                .font(.system(size: 64, weight: .bold))
                .foregroundStyle(scoreColor(score))
            Text(scoreLabel(score))
                .font(.title3)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity)
        .padding()
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 16))
    }

    private func nextEventCard(title: String, time: String) -> some View {
        VStack(spacing: 8) {
            Text("翌日の最初の予定")
                .font(.headline)
            Text(title)
                .font(.body)
                .foregroundStyle(.secondary)
            Text(time)
                .font(.title2.bold())
        }
        .frame(maxWidth: .infinity)
        .padding()
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 16))
    }

    private func focusPeakCard(start: Date?, end: Date?) -> some View {
        VStack(spacing: 8) {
            Text("集中ピーク時間帯")
                .font(.headline)
            if let start, let end {
                Text("\(timeString(start)) 〜 \(timeString(end))")
                    .font(.title2)
            } else {
                Text("—")
                    .font(.title2)
                    .foregroundStyle(.secondary)
            }
        }
        .frame(maxWidth: .infinity)
        .padding()
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 16))
    }

    private func bedtimeCard(bedtime: Date) -> some View {
        VStack(spacing: 8) {
            Text("推奨就寝時刻")
                .font(.headline)
            Text(timeString(bedtime))
                .font(.title2)
        }
        .frame(maxWidth: .infinity)
        .padding()
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 16))
    }

    private func sleepDebtCard(minutes: Int) -> some View {
        VStack(spacing: 8) {
            Text("睡眠負債")
                .font(.headline)
            Text(minutes > 0 ? "\(minutes) 分の借金あり" : "負債なし")
                .font(.title3)
                .foregroundStyle(minutes > 0 ? .orange : .green)
        }
        .frame(maxWidth: .infinity)
        .padding()
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 16))
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
