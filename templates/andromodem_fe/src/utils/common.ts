import type { FeatureAvailabilities } from "@/types/device"

export const getBaseUrl = (): string => {
    if (import.meta.env.MODE === "development") {
        // TODO: Move the development server URL to environment variables (.env) file
        return "http://localhost:49153"
    }
    return `${location.protocol}//${location.host}`
}

export const showModal = (modal_id: string) => {
    const modal = document.getElementById(modal_id) as HTMLDialogElement | null
    if (modal) {
        modal.showModal()
    }
}

export const isFeatureAvailable = (features: FeatureAvailabilities | null, key: string) => {
    if (features === null) {
        return false
    }
    return features?.some(
        (feature) => feature.key === key && feature.available === true
    )
}