import { describe, it, expect } from 'vitest';
import {
    parseDate,
    getTitle,
    getPages,
    getISBN,
    coverURLByISBN,
    computeYears,
    computeStats,
    computeBooksForYear,
    computeMonthlyBreakdown,
} from '../src/data.js';

const book = (over = {}) => ({
    Title: 'A Book',
    Author: 'An Author',
    'Number of Pages': 100,
    ISBN: '',
    ISBN13: '',
    'Date Read': '2025/03/15',
    Shelf: 'read',
    ...over,
});

describe('parseDate', () => {
    it('parses full YYYY/MM/DD', () => {
        expect(parseDate('2025/07/26')).toEqual({ year: 2025, month: 7, day: 26 });
    });

    it('parses year-only', () => {
        expect(parseDate('2024')).toEqual({ year: 2024, month: 0, day: 0 });
    });

    it('rejects empty, non-string, pre-1900, and out-of-range parts', () => {
        expect(parseDate('')).toBeNull();
        expect(parseDate(null)).toBeNull();
        expect(parseDate('1899/01/01')).toBeNull();
        expect(parseDate('2025/13/01')).toBeNull();
        expect(parseDate('2025/00/01')).toBeNull();
    });
});

describe('field accessors', () => {
    it('getTitle handles string and number', () => {
        expect(getTitle({ Title: '  Dune ' })).toBe('Dune');
        expect(getTitle({ Title: 1984 })).toBe('1984');
    });

    it('getPages handles number and numeric string', () => {
        expect(getPages({ 'Number of Pages': 320 })).toBe(320);
        expect(getPages({ 'Number of Pages': '448' })).toBe(448);
        expect(getPages({ 'Number of Pages': '' })).toBe(0);
    });

    it('getISBN prefers ISBN13 then ISBN', () => {
        expect(getISBN({ ISBN13: '9780593804216', ISBN: '059380421X' })).toBe('9780593804216');
        expect(getISBN({ ISBN13: '', ISBN: '059380421X' })).toBe('059380421X');
        expect(getISBN({ ISBN13: '', ISBN: '' })).toBe('');
    });

    it('coverURLByISBN builds Open Library URL or empty', () => {
        expect(coverURLByISBN('9780593804216')).toBe('https://covers.openlibrary.org/b/isbn/9780593804216-M.jpg');
        expect(coverURLByISBN('')).toBe('');
    });
});

describe('aggregations', () => {
    const books = [
        book({ 'Date Read': '2025/01/10', 'Number of Pages': 100 }),
        book({ 'Date Read': '2025/01/20', 'Number of Pages': 200 }),
        book({ 'Date Read': '2025/03/05', 'Number of Pages': 0 }),
        book({ 'Date Read': '2024/12/31', 'Number of Pages': 50 }),
        book({ 'Date Read': '2025/06/01', Shelf: 'to-read' }),
        book({ 'Date Read': '', Shelf: 'read' }),
    ];

    it('computeYears counts read books with valid dates', () => {
        const years = computeYears(books).sort((a, b) => b.year - a.year);
        expect(years).toEqual([
            { year: 2025, count: 3 },
            { year: 2024, count: 1 },
        ]);
    });

    it('computeStats totals pages and averages per month', () => {
        const stats = computeStats(books, 2025);
        expect(stats.totalBooks).toBe(3);
        expect(stats.totalPages).toBe(300);
        expect(stats.averagePerMonth).toBeCloseTo(3 / 12);
        expect(stats.monthlyBreakdown).toHaveLength(12);
    });

    it('computeMonthlyBreakdown places books in the right month', () => {
        const breakdown = computeMonthlyBreakdown(books.filter((b) => b.Shelf === 'read'));
        expect(breakdown[0].Count).toBe(2);
        expect(breakdown[2].Count).toBe(1);
        expect(breakdown[0].MonthName).toBe('Jan');
    });

    it('computeBooksForYear maps view models and filters by shelf', () => {
        const list = computeBooksForYear(books, 2025, 'read');
        expect(list).toHaveLength(3);
        expect(list[0]).toMatchObject({ title: 'A Book', author: 'An Author', shelf: 'read' });
        expect(list[0].month).toBe(1);
    });
});
