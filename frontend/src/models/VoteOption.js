class VoteOption {
  constructor(value, color) {
    this.value = value;
    this.color = color;
  }

  getColor() {
    return this.color;
  }

  static getAllOptions() {
    return [
      new VoteOption('1', '#3498db'),  // Light blue
      new VoteOption('2', '#2980b9'),  // Medium blue
      new VoteOption('3', '#16a085'),  // Teal
      new VoteOption('5', '#27ae60'),  // Green
      new VoteOption('8', '#f39c12'),  // Orange
      new VoteOption('13', '#e67e22'), // Dark orange
      new VoteOption('21', '#e74c3c'), // Red
      new VoteOption('34', '#9b59b6'), // Purple
      new VoteOption('?', '#95a5a6'),  // Gray
    ];
  }

  static getByValue(value) {
    return this.getAllOptions().find(option => option.value === value);
  }

  static getAllValues() {
    return this.getAllOptions().map(option => option.value);
  }

  static getAllNumericValues() {
    return this.getAllOptions().filter(option => option.value !== '?').map(option => option.value);
  }
}

export default VoteOption;
