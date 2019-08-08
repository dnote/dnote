// NoteData represents a data for a note as returned by services.
// The response from services need to conform to this interface.
export interface NoteData {
  uuid: string;
  created_at: string;
  updated_at: string;
  content: string;
  added_on: number;
  public: boolean;
  usn: number;
  book: {
    uuid: string;
    label: string;
  };
  user: {
    name: string;
    uuid: string;
  };
}
