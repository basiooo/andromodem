import type { Config } from "@/types/config"
import { getBaseUrl } from "@/utils/common"

console.log(getBaseUrl())
export const config: Config = {
    BASE_URL: getBaseUrl(),
    VERSION: import.meta.env.VITE_ANDROMODEM_VERSION ?? "unknown"
}