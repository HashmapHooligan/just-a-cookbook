import { defineBoot } from '#q-app/wrappers';
import { createI18n } from 'vue-i18n';

import messages from 'src/i18n';

export type MessageLanguages = keyof typeof messages;
export type MessageSchema = (typeof messages)['en-US'];

declare module 'vue-i18n' {
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  export interface DefineLocaleMessage extends MessageSchema {}
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  export interface DefineDateTimeFormat {}
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  export interface DefineNumberFormat {}
}

const savedLocale = (localStorage.getItem('locale') ?? 'en-US') as MessageLanguages;

export const i18n = createI18n<{ message: MessageSchema }, MessageLanguages>({
  locale: savedLocale,
  legacy: false,
  messages,
});

export default defineBoot(({ app }) => {
  app.use(i18n);
});
