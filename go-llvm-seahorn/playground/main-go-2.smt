(set-info :original "main-2-manual.ll")
(set-info :authors "SeaHorn v.10.0.0-rc0")
(declare-rel verifier.error (Bool Bool Bool ))
(declare-rel main@entry (Bool ))
(declare-rel main@block_false (Bool ))
(declare-rel main@block_true (Bool ))
(declare-rel main@block_true.split (Bool ))
(declare-var main@%phi.1_0 Int )
(declare-var main@%sm2_0 (Array Int Int) )
(declare-var main@%sm3_0 (Array Int Int) )
(declare-var main@%a.0_0 Int )
(declare-var main@%b.0_0 Int )
(declare-var main@%sm_0 (Array Int Int) )
(declare-var main@%sm1_0 (Array Int Int) )
(declare-var main@%a.ld.0_0 Int )
(declare-var main@%icmp.0_0 Bool )
(declare-var main@%zext.0_0 Int )
(declare-var main@%b.ld.0_0 Int )
(declare-var main@%icmp.1_0 Bool )
(declare-var main@%zext.1_0 Int )
(declare-var main@%a.ld.1_0 Int )
(declare-var main@%add.0_0 Int )
(declare-var main@%a.ld.2_0 Int )
(declare-var main@%b.ld.1_0 Int )
(declare-var main@%add.1_0 Int )
(declare-var main@%icmp.2_0 Bool )
(declare-var error.flag_0 Bool )
(declare-var error.flag_1 Bool )
(rule (verifier.error false false false))
(rule (verifier.error false true true))
(rule (verifier.error true false true))
(rule (verifier.error true true true))
(rule (=> (not error.flag_0) (main@entry error.flag_0)))
(rule (=> (and (main@entry error.flag_0)
         (not error.flag_0)
         (> main@%a.0_0 0)
         (> main@%b.0_0 0)
         (= main@%sm_0 (store main@%sm2_0 main@%a.0_0 90))
         (= main@%sm1_0 (store main@%sm3_0 main@%b.0_0 100))
         (= main@%a.ld.0_0 (select main@%sm_0 main@%a.0_0))
         (= main@%icmp.0_0 (> main@%a.ld.0_0 0))
         (= main@%zext.0_0 (ite main@%icmp.0_0 1 0))
         (or error.flag_0 main@%icmp.0_0)
         (= main@%b.ld.0_0 (select main@%sm1_0 main@%b.0_0))
         (= main@%icmp.1_0 (> main@%b.ld.0_0 0))
         (= main@%zext.1_0 (ite main@%icmp.1_0 1 0))
         (or error.flag_0 main@%icmp.1_0)
         (= main@%a.ld.1_0 (select main@%sm_0 main@%a.0_0))
         (= main@%add.0_0 (+ main@%a.ld.1_0 50))
         (= main@%a.ld.2_0 (select main@%sm_0 main@%a.0_0))
         (= main@%b.ld.1_0 (select main@%sm1_0 main@%b.0_0))
         (= main@%add.1_0 (+ main@%a.ld.2_0 main@%b.ld.1_0))
         (= main@%icmp.2_0 (> main@%add.0_0 main@%add.1_0))
         (or error.flag_0 main@%icmp.2_0))
    (main@block_true error.flag_0)))
(rule (=> (and (main@entry error.flag_0)
         (not error.flag_0)
         (> main@%a.0_0 0)
         (> main@%b.0_0 0)
         (= main@%sm_0 (store main@%sm2_0 main@%a.0_0 90))
         (= main@%sm1_0 (store main@%sm3_0 main@%b.0_0 100))
         (= main@%a.ld.0_0 (select main@%sm_0 main@%a.0_0))
         (= main@%icmp.0_0 (> main@%a.ld.0_0 0))
         (= main@%zext.0_0 (ite main@%icmp.0_0 1 0))
         (or error.flag_0 main@%icmp.0_0)
         (= main@%b.ld.0_0 (select main@%sm1_0 main@%b.0_0))
         (= main@%icmp.1_0 (> main@%b.ld.0_0 0))
         (= main@%zext.1_0 (ite main@%icmp.1_0 1 0))
         (or error.flag_0 main@%icmp.1_0)
         (= main@%a.ld.1_0 (select main@%sm_0 main@%a.0_0))
         (= main@%add.0_0 (+ main@%a.ld.1_0 50))
         (= main@%a.ld.2_0 (select main@%sm_0 main@%a.0_0))
         (= main@%b.ld.1_0 (select main@%sm1_0 main@%b.0_0))
         (= main@%add.1_0 (+ main@%a.ld.2_0 main@%b.ld.1_0))
         (= main@%icmp.2_0 (> main@%add.0_0 main@%add.1_0))
         (or error.flag_0 (not main@%icmp.2_0))
         (verifier.error true error.flag_0 error.flag_1))
    (main@block_false error.flag_1)))
(rule (=> (and (main@block_true error.flag_0) (not error.flag_0) (= main@%phi.1_0 1))
    (main@block_true.split error.flag_0)))
(rule (=> (and (main@entry error.flag_0) error.flag_0)
    (main@block_true.split error.flag_0)))
(rule (=> (and (main@block_false error.flag_0) error.flag_0)
    (main@block_true.split error.flag_0)))
(rule (=> (and (main@block_true error.flag_0) error.flag_0)
    (main@block_true.split error.flag_0)))
(query (main@block_true.split true))
