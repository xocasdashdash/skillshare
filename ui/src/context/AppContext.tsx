import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { api } from '../api/client';

interface AppContextValue {
  isProjectMode: boolean;
  projectRoot?: string;
}

const AppContext = createContext<AppContextValue>({ isProjectMode: false });

export function useAppContext() {
  return useContext(AppContext);
}

export function AppProvider({ children }: { children: ReactNode }) {
  const [value, setValue] = useState<AppContextValue>({ isProjectMode: false });

  useEffect(() => {
    api.getOverview().then((data) => {
      setValue({
        isProjectMode: data.isProjectMode,
        projectRoot: data.projectRoot,
      });
    }).catch(() => {
      // Keep defaults on error
    });
  }, []);

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}
