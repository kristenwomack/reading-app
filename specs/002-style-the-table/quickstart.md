# Quickstart: Validate Summary Table Layout

**Feature**: 002-style-the-table  
**Purpose**: Manual validation steps to verify the summary table displays correctly in 25%/75% layout.  
**Prerequisites**: Backend running on port 8080, frontend served, books.json populated with 2025 data

---

## Setup

1. **Start the backend server**:
   ```bash
   cd backend
   go run main.go
   ```
   Expected: Server starts on `http://localhost:8080`

2. **Access the application**:
   Open browser to `http://localhost:8080`

3. **Ensure test data exists**:
   Verify `books.json` contains at least one book with `finish_date` in 2025.

---

## Validation Steps

### ✅ Test 1: Desktop Layout (25%/75% Split)

**Scenario**: View dashboard on desktop screen (>768px width)

**Steps**:
1. Open browser to `http://localhost:8080`
2. Set browser window to 1440px width (or wider)
3. Select year "2025" from year selector

**Expected Results**:
- ✅ Summary card appears on LEFT side
- ✅ Summary card occupies approximately 25% of screen width
- ✅ Charts section appears on RIGHT side
- ✅ Charts section occupies approximately 75% of screen width
- ✅ Summary card and charts are SIDE-BY-SIDE (not stacked)
- ✅ Gap between summary and charts is visible (approximately 1.5rem/24px)

**How to Verify Width**:
- Open browser DevTools (F12)
- Inspect `.summary-section` element
- Check computed width ≈ 25% of viewport (accounting for gap)
- Inspect `.charts-section` element  
- Check computed width ≈ 75% of viewport (accounting for gap)

---

### ✅ Test 2: Summary Card Content

**Scenario**: Verify summary card displays correct statistics

**Steps**:
1. Ensure year 2025 is selected
2. Observe summary card content

**Expected Results**:
- ✅ Card title shows "2025 Summary" (or similar)
- ✅ "Total Books" row displays count (e.g., "42")
- ✅ "Total Pages" row displays sum with comma separator (e.g., "12,450")
- ✅ "Avg Pages/Book" row displays average (e.g., "296")
- ✅ All values match backend API response from `/api/stats?year=2025`

**How to Verify API Match**:
- Open browser DevTools Network tab
- Locate `GET /api/stats?year=2025` request
- Compare response JSON values to displayed values:
  ```json
  {
    "year": 2025,
    "total_books": 42,
    "total_pages": 12450,
    "avg_pages_per_book": 296
  }
  ```

---

### ✅ Test 3: Mobile Layout (Stacked)

**Scenario**: View dashboard on mobile screen (<768px width)

**Steps**:
1. Resize browser window to 375px width (iPhone SE size)
2. OR: Open DevTools → Toggle device toolbar → Select "iPhone SE"
3. Ensure year 2025 is selected

**Expected Results**:
- ✅ Summary card appears at TOP (full width)
- ✅ Charts section appears BELOW summary card (full width)
- ✅ Both sections are STACKED vertically (not side-by-side)
- ✅ Summary card width = 100% of viewport
- ✅ No horizontal scrolling required

**How to Verify**:
- Inspect `.dashboard-layout` element
- Check computed `grid-template-columns` is `1fr` (single column) or `flex-direction: column`

---

### ✅ Test 4: Responsive Breakpoint Transition

**Scenario**: Verify layout transitions smoothly at 768px breakpoint

**Steps**:
1. Start with browser width at 1440px (desktop layout)
2. Slowly drag browser window edge to reduce width
3. Observe layout changes as width crosses 768px threshold

**Expected Results**:
- ✅ Above 768px: Side-by-side layout maintained
- ✅ At/below 768px: Layout switches to stacked
- ✅ Transition is smooth (no content jumping or overflow)
- ✅ No horizontal scrollbars at any width

---

### ✅ Test 5: Empty State (No Books)

**Scenario**: View summary card when selected year has no books

**Steps**:
1. Select a year with no books (e.g., 2020 if no books tracked)
2. OR: Temporarily modify books.json to remove all 2025 entries

**Expected Results**:
- ✅ Summary card still displays
- ✅ Card shows "0 books" or "No books tracked for this year"
- ✅ No error messages or broken layout
- ✅ Charts section shows appropriate empty state (existing behavior)

---

### ✅ Test 6: Multiple Years

**Scenario**: Verify summary updates when switching years

**Steps**:
1. Select year 2025
2. Note summary statistics
3. Select different year (e.g., 2024)
4. Observe summary card updates

**Expected Results**:
- ✅ Summary card title updates to new year (e.g., "2024 Summary")
- ✅ Statistics update to reflect new year's data
- ✅ Layout remains intact (25%/75% on desktop)
- ✅ No flickering or layout shifts during update

---

### ✅ Test 7: Chart Rendering Integrity

**Scenario**: Verify existing charts still render correctly with new layout

**Steps**:
1. View dashboard with year 2025 selected (desktop width)
2. Observe charts in right section (75% width area)

**Expected Results**:
- ✅ All existing charts render without distortion
- ✅ Chart dimensions fit within 75% width area
- ✅ Chart legends and labels are readable
- ✅ No chart overflow or horizontal scrolling

**Note**: If charts appear stretched or compressed, the implementation may need to call `chart.resize()` after layout change.

---

## Cross-Browser Testing

Repeat Tests 1-3 on:
- ✅ Chrome (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ✅ Edge (latest)

**Expected**: Consistent layout across all browsers (CSS Grid has universal modern browser support)

---

## Performance Validation

**Metric**: Layout render time

**Steps**:
1. Open browser DevTools Performance tab
2. Start recording
3. Select a year to trigger layout update
4. Stop recording
5. Inspect "Layout" and "Paint" events

**Expected Results**:
- ✅ Total layout time < 16ms (60fps threshold)
- ✅ No forced synchronous layout warnings
- ✅ Paint events are minimal

---

## Failure Criteria

**FAIL if any of these occur**:
- ❌ Summary card width is not approximately 25% on desktop
- ❌ Layout does not stack on mobile (<768px)
- ❌ Summary statistics do not match API response
- ❌ Horizontal scrolling required at any width
- ❌ Charts are distorted or overlap with summary card
- ❌ Empty state shows error instead of "0 books"

---

## Success Criteria Summary

**All validation steps pass** = Feature is ready for release:
- Desktop layout: 25%/75% split ✅
- Mobile layout: Stacked ✅
- Summary card content accurate ✅
- Responsive transition smooth ✅
- Empty state handled gracefully ✅
- Year switching works ✅
- Charts render correctly ✅
- Cross-browser compatible ✅
- Performance acceptable ✅

---

## Rollback Plan

If validation fails:
1. Revert CSS changes in `frontend/styles/main.css`
2. Revert HTML changes in `frontend/index.html`
3. Revert JS changes in `frontend/src/ui.js`
4. Clear browser cache and hard reload
5. Verify original layout restored

---

**Last Updated**: 2025-01-16  
**Estimated Validation Time**: 15-20 minutes
