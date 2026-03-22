"use client";

import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { useI18n } from "../lib/i18n";

export default function ThemeLangToggle() {
  const [mounted, setMounted] = useState(false);
  const { theme, setTheme } = useTheme();
  const { lang, setLang } = useI18n();

  useEffect(() => {
    // Timeout prevents React hydration mismatch by delaying the state update
    const timeout = setTimeout(() => setMounted(true), 0);
    return () => clearTimeout(timeout);
  }, []);

  if (!mounted) {
    return (
      <div className="flex items-center gap-2">
        <div className="w-8 h-8 rounded-lg bg-[var(--color-card)] animate-pulse" />
        <div className="w-16 h-8 rounded-lg bg-[var(--color-card)] animate-pulse" />
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <button
        onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        className="p-2 rounded-lg bg-[var(--color-card)] border border-[var(--color-border)] text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-colors"
        title={
          theme === "dark" ? "Switch to light mode" : "Switch to dark mode"
        }
      >
        {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
      </button>

      <button
        onClick={() => setLang(lang === "en" ? "vi" : "en")}
        className="px-3 py-1.5 rounded-lg bg-[var(--color-card)] border border-[var(--color-border)] text-sm font-medium text-[var(--color-foreground)] hover:bg-[var(--color-card-hover)] transition-colors"
        title={lang === "en" ? "Chuyển sang tiếng Việt" : "Switch to English"}
      >
        {lang === "en" ? "EN" : "VI"}
      </button>
    </div>
  );
}
