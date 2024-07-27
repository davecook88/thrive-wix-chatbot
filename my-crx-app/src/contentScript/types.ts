export interface CustomerInformation {
  spanishLevel: 'Beginner' | 'Intermediate' | 'Advanced';
  notes: string;
  upcomingClasses: { date: string; topic: string }[];
}
