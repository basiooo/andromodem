import { type FC,useEffect, useState } from "react"

import { THEMES } from "@/constants/themes"

const ThemeSelector: FC = () => {

    const [selectedTheme, setSelectedTheme] = useState<string | null>(null)

    useEffect(() => {
        const theme = document.documentElement.getAttribute("data-theme")
        if (theme) {
            setSelectedTheme(theme)
        }
    }, [])

    return (
        <div className="flex gap-2">
            <div
                title="Change Theme"
                className="dropdown dropdown-end block"
            >
                <div
                    tabIndex={0}
                    role="button"
                    className="btn btn-sm gap-1"
                >
                    <div
                        className="bg-base-100 border-base-content/10 grid shrink-0 grid-cols-2 gap-0.5 rounded-md border p-1">
                        <div className="bg-base-content size-1 rounded-full"></div>
                        <div className="bg-primary size-1 rounded-full"></div>
                        <div className="bg-secondary size-1 rounded-full"></div>
                        <div className="bg-accent size-1 rounded-full"></div>
                    </div>
                    <svg
                        width="12px"
                        height="12px"
                        className="mt-px hidden h-2 w-2 fill-current opacity-60 sm:inline-block"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 2048 2048"
                    >
                        <path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"></path>
                    </svg>
                    <span className="hidden md:block">Change Theme</span>
                </div>
                <div
                    tabIndex={0}
                    className="dropdown-content bg-base-200 text-base-content rounded-box top-px h-[30.5rem] max-h-[calc(100vh-8.6rem)] overflow-y-auto border border-white/5 shadow-2xl outline-1 outline-black/5"
                >
                    <ul className="menu w-56">
                        <li className="menu-title text-xs">Change Theme</li>
                        {THEMES.map((theme) => (
                            <li key={theme}>
                                <button
                                    className="gap-3 px-2" data-set-theme={theme}
                                    onClick={() => setSelectedTheme(theme)}
                                >
                                    <div
                                        data-theme={theme}
                                        className="bg-base-100 grid shrink-0 grid-cols-2 gap-0.5 rounded-md p-1 shadow-sm"
                                    >
                                        <div className="bg-base-content size-1 rounded-full"></div>
                                        <div className="bg-primary size-1 rounded-full"></div>
                                        <div className="bg-secondary size-1 rounded-full"></div>
                                        <div className="bg-accent size-1 rounded-full"></div>
                                    </div>
                                    <div className="w-32 truncate">{theme}</div>
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        width="16"
                                        height="16"
                                        viewBox="0 0 24 24"
                                        fill="currentColor"
                                        className={`h-3 w-3 shrink-0 ${selectedTheme === theme ? "visible" : "invisible"
                                            }`}
                                    >
                                        <path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z" />
                                    </svg>
                                </button>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
    )
}

export default ThemeSelector