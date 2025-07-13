import axios from "axios"

import type { GitHubRelease } from "@/types/update"

const GITHUB_API_BASE = "https://api.github.com"
const REPO_OWNER = "basiooo"
const REPO_NAME = "andromodem"

export const updateApi = {
    getLatestRelease: async (): Promise<GitHubRelease> => {
        const { data } = await axios.get<GitHubRelease>(
            `${GITHUB_API_BASE}/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest`
        )
        return data
    },
    
    getAllReleases: async (): Promise<GitHubRelease[]> => {
        const { data } = await axios.get<GitHubRelease[]>(
            `${GITHUB_API_BASE}/repos/${REPO_OWNER}/${REPO_NAME}/releases`
        )
        return data
    }
}