// src/lib/stores/viewport.ts

import { writable } from 'svelte/store';

const createViewportStore = () => {
  const { subscribe, set } = writable({
    isMobile: false,
    width: 0,
    height: 0,
  });

  const checkViewport = () => {
    set({
      isMobile: window.innerWidth < 768,
      width: window.innerWidth,
      height: window.innerHeight,
    });
  };

  const cleanup = () => {
    if (typeof window !== 'undefined') {
      window.removeEventListener('resize', checkViewport);
    }
  };

  if (typeof window !== 'undefined') {
    checkViewport();
    window.addEventListener('resize', checkViewport);
  }

  return {
    subscribe,
    cleanup,
  };
};

export const viewport = createViewportStore();
