/*
1. Change the website eld, so it only contains the domain.
‣ Example: https://domain.com/index.php → domain.com
*/
UPDATE "MY_TABLE" SET website=substring(website from '^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)');

/*
2. Count how many spots contain the same domain.
*/
SELECT website AS "domain", count(id) AS "numberOfspots" FROM "MY_TABLE" GROUP BY website HAVING website  IS NOT NULL;


/*
3. Return spots which have a domain with a count greater than 1
*/
SELECT * FROM "MY_TABLE" mt
WHERE website IN (SELECT website FROM "MY_TABLE" mt GROUP BY website HAVING count(id) > 1 AND website IS NOT NULL);

/*
4. Make a PL/SQL function for point 1 above.
//test ==> select  replaceWebsiteToDomain('^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)');
*/
create  OR REPLACE function replaceWebsiteToDomain(regex varchar)
returns table (
  spotName varchar,
  website varchar
)
as
$$
begin
 	-- logic
UPDATE "MY_TABLE" SET "website" = substring(public."MY_TABLE"."website" from regex);
return query
select public."MY_TABLE"."name"  as "spotName", public."MY_TABLE"."website" from public."MY_TABLE";
end;
$$
LANGUAGE plpgsql;