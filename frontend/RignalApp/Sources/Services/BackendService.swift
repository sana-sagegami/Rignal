import Foundation

struct BackendService {
    private let baseURL: String

    init() {
        guard let url = Bundle.main.infoDictionary?["BackendBaseURL"] as? String else {
            fatalError("BackendBaseURL not set in Info.plist")
        }
        self.baseURL = url
    }

    func fetchSummary(nextEvent: String? = nil) async throws -> SummaryResponse {
        var components = URLComponents(string: "\(baseURL)/summary")!
        if let nextEvent {
            components.queryItems = [URLQueryItem(name: "next_event", value: nextEvent)]
        }
        guard let url = components.url else { throw URLError(.badURL) }

        let (data, _) = try await URLSession.shared.data(from: url)

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .iso8601
        return try decoder.decode(SummaryResponse.self, from: data)
    }
}
