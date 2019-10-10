/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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

export interface EmailPrefData {
  digestWeekly: boolean;
}

export interface UserData {
  uuid: string;
  email: string;
  emailVerified: boolean;
  pro: boolean;
  classic: boolean;
}

export type BookData = {
  uuid: string;
  usn: number;
  created_at: string;
  updated_at: string;
  label: string;
};

// BookDomain is the possible values for the field in the repetition_rule
// indicating how to derive the source books for the repetition_rule.
export enum BookDomain {
  // All incidates that all books are eligible to be the source books
  All = 'all',
  // Including incidates that some specified books are eligible to be the source books
  Including = 'including',
  // Excluding incidates that all books except for some specified books are eligible to be the source books
  Excluding = 'excluding'
}

export interface RepetitionRuleData {
  uuid: string;
  title: string;
  enabled: boolean;
  hour: number;
  minute: number;
  bookDomain: BookDomain;
  frequency: number;
  books: BookData[];
  lastActive: number;
  noteCount: number;
  createdAt: string;
  updatedAt: string;
}
