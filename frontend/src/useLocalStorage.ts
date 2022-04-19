import { Dispatch, SetStateAction, useEffect, useState } from 'react';

function useLocalStorage<T>(storageKey: string, fallbackState: any): [T, Dispatch<SetStateAction<T>>] {
  const [value, setValue] = useState<T>(
    JSON.parse(localStorage.getItem(storageKey)!) ?? fallbackState
  );

  useEffect(() => {
    localStorage.setItem(storageKey, JSON.stringify(value));
  }, [value, storageKey]);

  return [value, setValue];
}

export default useLocalStorage;