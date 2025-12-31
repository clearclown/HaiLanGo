import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#4A90E2',
          light: '#6BA5E9',
          dark: '#3A7BC8',
        },
        secondary: {
          DEFAULT: '#50C878',
          light: '#6DD58E',
          dark: '#3FB066',
        },
        accent: {
          DEFAULT: '#FF6B6B',
          light: '#FF8585',
          dark: '#E55555',
        },
        background: {
          DEFAULT: '#FFFFFF',
          secondary: '#F5F7FA',
          tertiary: '#EEF2F7',
        },
        text: {
          primary: '#2C3E50',
          secondary: '#7F8C8D',
          muted: '#A0AEC0',
        },
        border: {
          DEFAULT: '#E0E6ED',
          light: '#F0F4F8',
          dark: '#CBD5E1',
        },
        success: {
          DEFAULT: '#27AE60',
          light: '#D4EDDA',
          dark: '#1E8449',
        },
        warning: {
          DEFAULT: '#F39C12',
          light: '#FFF3CD',
          dark: '#D68910',
        },
        error: {
          DEFAULT: '#E74C3C',
          light: '#F8D7DA',
          dark: '#C0392B',
        },
        info: {
          DEFAULT: '#3498DB',
          light: '#D1ECF1',
          dark: '#2980B9',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Noto Sans JP', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        soft: '0 2px 8px rgba(0, 0, 0, 0.06)',
        medium: '0 4px 16px rgba(0, 0, 0, 0.08)',
        hard: '0 8px 24px rgba(0, 0, 0, 0.12)',
      },
      animation: {
        'fade-in': 'fadeIn 0.2s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'pulse-slow': 'pulse 3s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
};

export default config;
