'use client';

import { animate } from 'framer-motion';
import { useEffect, useState } from 'react';

const delimiter = ''; // or " " to split by word

export function useAnimatedText(text: string,time :number) {
  const [cursor, setCursor] = useState(0);
  const [startingCursor, setStartingCursor] = useState(0);
  const [prevText, setPrevText] = useState(text);

  if (prevText !== text && text) {
    setPrevText(text);
    setStartingCursor(text.startsWith(prevText) ? cursor : 0);
  }

  useEffect(() => {
    if(text){
      const controls = animate(startingCursor, text?.split(delimiter).length, {
        // Tweak the animation here
        duration: time,
        ease: 'easeOut',
        onUpdate(latest) {
          setCursor(Math.floor(latest));
        },
      });

      return () => controls.stop();
    }
  }, [startingCursor, text]);

  return {
    text:  text ?text?.split(delimiter)?.slice(0, cursor)?.join(delimiter) : '',
    done: cursor === text?.split(delimiter)?.length,
  };
}
