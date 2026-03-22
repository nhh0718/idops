"use client";

import { createContext, useContext, useEffect, useState } from "react";
import en from "./dictionaries/en.json";
import vi from "./dictionaries/vi.json";

type Language = "en" | "vi";
type Dictionary = typeof en;

interface I18nContextType {
  lang: Language;
  setLang: (lang: Language) => void;
  t: (key: string) => string;
}

const I18nContext = createContext<I18nContextType | null>(null);

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [lang, setLangState] = useState<Language>("en");
  const [dict, setDict] = useState<Dictionary>(en);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const savedLang = localStorage.getItem("idops-lang") as Language;

    // Use timeout to avoid direct state mutation during hydration phase
    const timeout = setTimeout(() => {
      if (savedLang && (savedLang === "en" || savedLang === "vi")) {
        setLangState(savedLang);
        setDict(savedLang === "en" ? en : vi);
      }
      setMounted(true);
    }, 0);

    return () => clearTimeout(timeout);
  }, []);

  const setLang = (newLang: Language) => {
    setLangState(newLang);
    setDict(newLang === "en" ? en : vi);
    localStorage.setItem("idops-lang", newLang);
  };

  const t = (path: string): string => {
    if (!mounted) return "";

    const keys = path.split(".");
    let current: Record<string, unknown> = dict;

    for (const key of keys) {
      if (current === undefined || current[key] === undefined) {
        console.warn(`Translation key not found: ${path}`);
        return path;
      }
      current = current[key] as Record<string, unknown>;
    }

    return current as unknown as string;
  };

  return (
    <I18nContext.Provider value={{ lang, setLang, t }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useI18n() {
  const context = useContext(I18nContext);
  if (!context) {
    throw new Error("useI18n must be used within an I18nProvider");
  }
  return context;
}
