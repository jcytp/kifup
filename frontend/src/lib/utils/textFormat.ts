export const formatDateTime = (dateStr: string | Date): string => {
  if (!dateStr) {
    return '';
  }
  const date = new Date(dateStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

export const formatTimeRule = (initial: string, byoyomi: string, increment: string): string => {
  let result = '';
  if (initial) {
    const v = parseInt(initial.replace('秒', ''));
    if (v) {
      result += Math.floor(v / 60).toString() + '分';
    }
  }
  if (byoyomi) {
    const v = parseInt(byoyomi.replace('秒', ''));
    if (v) {
      if (result != '') {
        result += ', ';
      }
      result += '秒読み ' + v.toString() + '秒';
    }
  }
  if (increment) {
    const v = parseInt(byoyomi.replace('秒', ''));
    if (v) {
      if (result != '') {
        result += ', ';
      }
      result += '加算 ' + v.toString() + '秒';
    }
  }
  return result;
};
