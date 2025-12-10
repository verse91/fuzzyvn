# Đóng góp cho FuzzyVN

Cảm ơn bạn đã quan tâm đến FuzzyVN. Chúng tôi rất hoan nghênh những đóng góp từ cộng đồng và trân trọng thời gian cũng như công sức bạn bỏ ra để giúp dự án tốt hơn. Trước khi bắt đầu, vui lòng dành chút thời gian đọc qua các hướng dẫn này.

Cần liên hệ? Hãy nhắn tin cho tôi qua Telegram [@justtheverse](t.me/JustTheVerse).

## 1 Cách thức Đóng góp

Chúng tôi hoan nghênh mọi hình thức đóng góp, bao gồm:

  - Sửa lỗi (Bug fixes)
  - Cải thiện tài liệu
  - Viết test và tối ưu hóa hiệu năng

Để đóng góp, hãy làm theo các bước sau:

  - **Fork** repository và tạo một branch mới.
  - **Thực hiện thay đổi (Make your changes)**, đảm bảo tuân thủ các tiêu chuẩn lập trình của chúng tôi (xem mục §3 bên dưới).
  - **Chạy tests** để đảm bảo các thay đổi của bạn không làm hỏng các chức năng hiện có.
  - **Gửi pull request (PR)** kèm theo mô tả rõ ràng về những thay đổi.
  - Tôi hoặc một maintainer khác sẽ xem xét PR của bạn, đề xuất thay đổi nếu cần thiết và merge (hợp nhất) khi đã được chấp thuận.

### Commit

Vui lòng sử dụng nội dung commit rõ ràng và mô tả đúng trọng tâm. Phải sử dụng tiếng anh theo quy chuẩn của chúng tôi. Hãy tuân theo định dạng "conventional commit" khi có thể:

  + `feat(fuzzyvn): enhance this function`
  + `fix(test): fix benchmark test`
  + `docs(root): update README with new set up steps`

## 2. Điều khoản Pháp lý

Bằng việc gửi đóng góp, bạn tuyên bố và đảm bảo rằng:

  - Đây là công việc gốc của bạn, hoặc bạn có đủ quyền hạn để gửi nó.
  - Bạn cấp cho các maintainer (người duy trì) và người dùng FuzzyVN quyền sử dụng, sửa đổi và phân phối nó theo giấy phép 0BSD (xem file LICENSE).
  - Trong phạm vi đóng góp của bạn có liên quan đến bằng sáng chế, bạn cấp giấy phép vĩnh viễn, toàn cầu, không độc quyền, miễn phí bản quyền và không thể thu hồi cho các maintainer và người dùng FuzzyVN để chế tạo, sử dụng, bán, chào bán, nhập khẩu và chuyển giao đóng góp của bạn như một phần của dự án.

Chúng tôi không yêu cầu Thỏa thuận Cấp phép Người đóng góp (CLA). Tuy nhiên, bằng việc đóng góp, bạn đồng ý cấp phép cho nội dung của mình theo các điều khoản tương thích với Giấy phép 0BSD và cấp các quyền sáng chế được mô tả ở trên. Nếu đóng góp của bạn bao gồm mã nguồn của bên thứ ba, bạn chịu trách nhiệm đảm bảo nó tương thích với 0BSD và được ghi nhận nguồn gốc hợp lệ.

Ở những nơi luật pháp cho phép, bạn từ bỏ mọi quyền nhân thân đối với đóng góp của mình (ví dụ: quyền phản đối việc sửa đổi). Nếu các quyền đó không thể từ bỏ, bạn đồng ý không thực thi chúng theo cách gây cản trở việc dự án sử dụng đóng góp của bạn.

## 3. Tiêu chuẩn Lập trình

Để duy trì sự nhất quán cho mã nguồn, vui lòng tuân thủ các nguyên tắc sau:

  - Sử dụng phong cách và quy ước lập trình hiện có của dự án.
  - Đảm bảo mọi thay đổi code đều được chú thích và tài liệu hóa tốt.
  - Viết test cho các tính năng mới và các bản sửa lỗi.
  - Tránh thêm các dependency (thư viện phụ thuộc) không cần thiết.

## 4. Báo cáo Vấn đề (Reporting Issues)

Nếu bạn tìm thấy lỗi hoặc có yêu cầu tính năng mới, vui lòng mở một Issue và cung cấp càng nhiều chi tiết càng tốt:

  - Các bước để tái hiện, bao gồm hệ điều hành và phiên bản FuzzyVN.
  - Hành vi mong đợi và hành vi thực tế.
  - Nguyên nhân nghi ngờ (nếu có).

## 5. Ghi nhận

Chúng tôi sử dụng tiêu chuẩn All Contributors để ghi nhận các thành viên cộng đồng. Nếu đóng góp của bạn được merge, bạn sẽ được thêm vào danh sách người đóng góp của dự án. Điều này bao gồm mọi hình thức đóng góp—từ code, tài liệu, thiết kế, cho đến kiểm thử (testing) và nhiều hơn nữa.