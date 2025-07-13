import axios from "axios"

import { config } from "@/config"

export const apiClient = axios.create({
    baseURL: `${config.BASE_URL}/api/`,
    headers: {
        "Content-type": "application/json"
    }
})