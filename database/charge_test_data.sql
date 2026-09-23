-- 收费模块接口测试数据。
-- 执行后使用 POST /api/charge/getPendingCharges，body: {"encounter_id": <查询到的 ID>}。
-- 本脚本只重置 patient_id = TEST_CHARGE_001 的测试就诊及其收费记录。

USE his_emergency;

INSERT INTO emergency_encounter
    (visit_no, patient_id, patient_name, patient_age, gender, phone, charge_type,
     encounter_type, dept_id, dept_name, doctor_id, doctor_name, encounter_status,
     total_amount, paid_amount, payment_status, remark)
VALUES
    ('TESTCHARGE20260922', 'TEST_CHARGE_001', '收费接口测试患者', '30岁', '男', '13800000001', '自费',
     'OPD', 'EMG', '急诊科', 'DOC_TEST_001', '测试医生', 'REGISTERED', 0, 0, 'UNPAID', '收费模块测试数据')
ON DUPLICATE KEY UPDATE
    patient_name = VALUES(patient_name), encounter_status = 'REGISTERED',
    total_amount = 0, paid_amount = 0, payment_status = 'UNPAID', remark = VALUES(remark);

SET @encounter_id = (SELECT id FROM emergency_encounter WHERE visit_no = 'TESTCHARGE20260922');

-- 清除这条测试就诊此前产生的收费记录，使待收费接口每次均可验证。
DELETE FROM emergency_charge WHERE encounter_id = @encounter_id;

INSERT INTO emergency_prescription
    (prescription_no, encounter_id, patient_id, patient_name, presc_type, status,
     prescriber_id, prescriber_name, total_amount, drug_amount, visit_no, visit_date, memo)
VALUES
    ('TESTRX20260922001', @encounter_id, 'TEST_CHARGE_001', '收费接口测试患者', 'WESTERN', 'SIGNED',
     'DOC_TEST_001', '测试医生', 25.50, 25.50, 'TESTCHARGE20260922', CURDATE(), '收费模块测试处方')
ON DUPLICATE KEY UPDATE
    encounter_id = VALUES(encounter_id), status = 'SIGNED', total_amount = 25.50, drug_amount = 25.50;

SET @prescription_id = (SELECT id FROM emergency_prescription WHERE prescription_no = 'TESTRX20260922001');

DELETE FROM emergency_prescription_item WHERE prescription_id = @prescription_id;

INSERT INTO emergency_prescription_item
    (prescription_id, item_no, item_class, drug_code, drug_name, drug_spec,
     total_qty, total_unit, unit_price, amount, dispense_status, remark)
VALUES
    (@prescription_id, 1, 'A', 'TEST_DRUG_001', '收费测试药品', '0.5g*20片',
     3, '盒', 8.50, 25.50, 'PENDING', '收费模块测试明细');

SELECT @encounter_id AS encounter_id, @prescription_id AS prescription_id;
