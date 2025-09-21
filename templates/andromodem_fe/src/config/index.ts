import type { Config } from "@/types/config"
import { getBaseUrl, getBaseUrlWs } from "@/utils/common"

console.log(getBaseUrl())
export const config: Config = {
    BASE_URL: getBaseUrl(),
    BASE_URL_WS: getBaseUrlWs(),
    VERSION: import.meta.env.VITE_ANDROMODEM_VERSION ?? "unknown"
}